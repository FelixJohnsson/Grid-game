import { Resources } from '../types';

const InformationTopBar = ({ resources }: { resources: Resources }) => {
  return (
    <div className="bg-gray-800 text-white p-4 flex justify-center">
      <div className="mr-8">Inhabitants: {resources.inhabitants.length}</div>
      <div className="mr-8">Wood: {resources.wood}</div>
      <div className="mr-8">Food: {resources.food}</div>
      <div className="mr-8">Stone: {resources.stone}</div>
      <div className="mr-8">Ores: {resources.ores}</div>
      <div className="mr-8">Money: {resources.money}</div>
    </div>
  );
};

export default InformationTopBar;